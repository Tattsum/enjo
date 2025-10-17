import { render, screen } from '@testing-library/react';
import LevelSlider from '../LevelSlider';

describe('LevelSlider', () => {
  it('renders slider with current value', () => {
    const onChange = jest.fn();
    render(<LevelSlider value={3} onChange={onChange} />);

    const slider = screen.getByRole('slider');
    expect(slider).toHaveValue('3');
  });

  it('displays the current level label', () => {
    const onChange = jest.fn();
    render(<LevelSlider value={3} onChange={onChange} />);

    expect(screen.getByText('レベル 3')).toBeInTheDocument();
  });

  it('calls onChange when slider is moved', async () => {
    const onChange = jest.fn();
    render(<LevelSlider value={3} onChange={onChange} />);

    const slider = screen.getByRole('slider');

    // Simulate changing the slider value
    const changeEvent = new Event('change', { bubbles: true });
    Object.defineProperty(slider, 'value', { value: '4', writable: true });
    slider.dispatchEvent(changeEvent);

    expect(onChange).toHaveBeenCalledWith(4);
  });

  it('has min value of 1', () => {
    const onChange = jest.fn();
    render(<LevelSlider value={1} onChange={onChange} />);

    const slider = screen.getByRole('slider') as HTMLInputElement;
    expect(slider.min).toBe('1');
  });

  it('has max value of 5', () => {
    const onChange = jest.fn();
    render(<LevelSlider value={5} onChange={onChange} />);

    const slider = screen.getByRole('slider') as HTMLInputElement;
    expect(slider.max).toBe('5');
  });

  it('displays level 1 description', () => {
    const onChange = jest.fn();
    render(<LevelSlider value={1} onChange={onChange} />);

    expect(screen.getByText(/少し配慮に欠ける表現/)).toBeInTheDocument();
  });

  it('displays level 3 description', () => {
    const onChange = jest.fn();
    render(<LevelSlider value={3} onChange={onChange} />);

    expect(screen.getByText(/明確に批判されそうな表現/)).toBeInTheDocument();
  });

  it('displays level 5 description', () => {
    const onChange = jest.fn();
    render(<LevelSlider value={5} onChange={onChange} />);

    expect(screen.getByText(/炎上確実な表現/)).toBeInTheDocument();
  });

  it('applies different color for each level', () => {
    const onChange = jest.fn();
    const { rerender } = render(<LevelSlider value={1} onChange={onChange} />);

    let slider = screen.getByRole('slider');
    expect(slider).toHaveClass('accent-blue-500');

    rerender(<LevelSlider value={5} onChange={onChange} />);
    slider = screen.getByRole('slider');
    expect(slider).toHaveClass('accent-fire-700');
  });

  it('renders with proper accessibility attributes', () => {
    const onChange = jest.fn();
    render(<LevelSlider value={3} onChange={onChange} />);

    const slider = screen.getByRole('slider');
    expect(slider).toHaveAttribute('aria-label', '炎上レベル');
  });
});
